public class S extends Struct_int_v implements SPtr{

	public S(int v){
		super(v);
	}

	public S(){
		super();
	}

	public int Square(){
		return Square(this);
	}

	public static int Square(Struct_int_v s){
		return s.v * s.v;
	}


	public void Inc(){
		Inc(this);
	}

	public static void Inc(Struct_int_v s){
		s.v++;
	}

	
}
