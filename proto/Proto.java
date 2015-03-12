//import go.;

public class Proto{
	

	static int gi;

	static go.Ptr_int gi_addr = new go.Ptr_int(){
		public int value(){return gi;}
		public void set(int v){gi = v;}
		public boolean equals(Object o){return this == o;}
	};

	public static void main(String[] args){

		// basic
		
		int i = 0;
		int i2 = 2;

		i = 7;
		i2 = i;

		System.out.println(i);
		System.out.println(i2);

		// escaped
		final go.Int e = new go.Int(3);

		i = e.value();
		e.set(i2);
		e.set(e.value() + 1);

		System.out.println(e.value());


		System.out.println("pointer->basic");
		go.Ptr_int eptr = new go.IntWrapper(e);
		System.out.println(eptr.value());

		go.Ptr_int giptr = gi_addr;
		go.Ptr_int giptr2 = gi_addr;
		giptr.set(3);
		System.out.println(giptr2.value());
		System.out.println(giptr.equals(giptr2));

		go.Ptr_int gi = gi_addr;
		System.out.println(gi.equals(gi_addr));
		
		
		System.out.println("struct");
		final Struct_int_v xs = new Struct_int_v();
		final Struct_int_v xs2 = new Struct_int_v(7);
		xs.set(xs2);
		
		System.out.println(xs.v);
		System.out.println(xs2.v);
		System.out.println(xs.equals(xs2));

		System.out.println("pointer->struct");
		Struct_int_v xsptr = new Struct_int_v();
		Struct_int_v xsptr2 = xs2;

		System.out.println(xsptr2 == xs2);
		System.out.println(xsptr == xs2);

		System.out.println(xsptr.value() == xs2.value());
		xsptr.set(xs2);
		System.out.println(xsptr.equals(xs2));


		System.out.println("named->basic");
		
		int mi = 0;

		System.out.println(mi);
		System.out.println(MyInt.Square(mi));

		System.out.println("named->struct");
		
		final Struct_int_v s = new Struct_int_v();
		System.out.println(s.v);
		s.set(xs);	
		System.out.println(s.v);
		xs.v = 99;
		xs.set(s);
		System.out.println(S.Square(s));

		final Struct_int_v t = new Struct_int_v();
		t.set(s);
		System.out.println(t.v);

		System.out.println("pointer->named->struct");
		Ptr_Struct_int_v sptr = new Ptr_Struct_int_v(s);
		Ptr_S.Inc(sptr);
		System.out.println(sptr.equals(new Ptr_Struct_int_v(s)));
		System.out.println(S.Square(s));

		sptr = new Ptr_Struct_int_v(xs);
		System.out.println(sptr.equals(new Ptr_Struct_int_v(xs)));
		sptr.set(new Struct_int_v(33));
		System.out.println(xs.v);

		System.out.println("interface");
		Object any = null;

		any = new go.Int(i);
		System.out.println(any.equals(new go.Int(3)));

		//int _1 = 0;
		//boolean ok=false;
		//if(any instanceof go.Int){
		//	_1 = ((go.Int)any).value;
		//	ok = true;
		//}
		//System.out.println(ok);

		//go.Ptr_int _2;
		//ok = false;
		//if(any instanceof go.Ptr_int){
		//	_2 = ((go.Ptr_int)any);
		//	ok = true;
		//}
		//System.out.println(ok);

		any = e;
		System.out.println(any.equals(new go.Int(3)));

		any = eptr;
		System.out.println(any.equals(new go.Int(3)));

		System.out.println("func");
		final go.Int ncalls = new go.Int();
		go.Func_int_int f = null;
		f = new go.Func_int_int(new go.Interface_call_int_int(){
			public int call(int x){
				ncalls.value++;
				return x*x;
			}
		});
		final go.Func_int_int a = new go.Func_int_int(f);
		System.out.println(a.call(3));
		System.out.println(ncalls.value);
		a.set(new go.Interface_call_int_int(){public int call(int x){return someF(x);}});
		
		System.out.println("named->func");
		go.Func_int_int myF = f;
		System.out.println(MyF.do_(myF, 4));

	}

	public static int someF(int x){
		return x*(x+1);
	}
}
