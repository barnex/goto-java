public class Struct_int_v{

	
	int v;
	
	public Struct_int_v(){}

	
	public Struct_int_v(Struct_int_v other){
		this.v = other.v;
	}
	
	public Struct_int_v(int v){
		this.v = v;
	}
	
	public void set(Struct_int_v  other){this.v = other.v;
	}
}
